export interface Store<Value> {
  get(): Value;
  set(value: Value): void;
  subscribe(listener: (value: Value, oldValue?: Value) => void): () => void;
  listen(listener: (value: Value, oldValue?: Value) => void): () => void;
  value: Value;
  off(): void;
  eager(): void;
  onMount(fn: () => void | (() => void)): void; // <-- Add this
  // Internal:
  _onMounts?: Array<() => void | (() => void)>;
  _cleanup?: () => void;
  _mounted?: boolean;
}

export function store<Value>(initialValue: Value): Store<Value> {
  let value = initialValue;
  let listeners: Array<(v: Value, o?: Value) => void> = [];
  const _onMounts: Array<() => void | (() => void)> = [];
  let _cleanup: (() => void) | undefined;
  let _mounted = false;

  function runOnMountIfNeeded() {
    if (!_mounted && _onMounts.length) {
      _mounted = true;
      for (const fn of _onMounts) {
        const cleanup = fn();
        if (typeof cleanup === "function") _cleanup = cleanup;
      }
    }
  }

  return {
    get: () => {
      runOnMountIfNeeded();
      return value;
    },
    set: (newValue: Value) => {
      if (value !== newValue) {
        const oldValue = value;
        value = newValue;
        for (const l of listeners) l(value, oldValue);
      }
    },
    subscribe: (listener) => {
      listeners.push(listener);
      runOnMountIfNeeded();
      listener(value);
      return () => {
        listeners = listeners.filter((l) => l !== listener);
        if (listeners.length === 0 && _cleanup) {
          _cleanup();
          _cleanup = undefined;
        }
      };
    },
    listen(listener) {
      return this.subscribe(listener);
    },
    get value() {
      runOnMountIfNeeded();
      return value;
    },
    off: () => {
      if (_cleanup) {
        _cleanup();
        _cleanup = undefined;
      }
      listeners = [];
    },
    eager: runOnMountIfNeeded,
    onMount(fn) {
      _onMounts.push(fn);
    },
    _onMounts,
    _cleanup,
    _mounted,
  };
}

