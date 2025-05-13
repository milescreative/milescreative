import { createFileRoute } from '@tanstack/react-router';
import { useState } from 'react';

// Define a type for the data expected from the token endpoint
interface CsrfTokenResponse {
  token: string;
}

// Define a type for the data expected from the protected API endpoint
interface ProtectedApiResponse {
  message: string;
  receivedData: any; // Adjust this type based on what your protected API actually returns
}

// --- csrf-manager.ts (simulated within the same file for simplicity) ---
// In a real application, this would be in a separate file (e.g., src/utils/csrfManager.ts)
let csrfToken: string | null = null; // Initialize with null and type as string or null

function getCsrfToken(): string | null {
  console.log("Getting CSRF token from internal storage.");
  return csrfToken;
}

function setCsrfToken(token: string): void { // Explicitly type the parameter and return
  console.log("Setting CSRF token in internal storage:", token);
  csrfToken = token;
}

function clearCsrfToken(): void { // Explicitly type the return
  console.log("Clearing CSRF token from internal storage.");
  csrfToken = null;
}
// --- End of csrf-manager.ts simulation ---

export const Route = createFileRoute('/csrf')({
  component: RouteComponent,
});

function RouteComponent() {
  const [apiResponse, setApiResponse] = useState<string>('...'); // Type the state

  const fetchCsrfToken = async () => {
    console.log("Fetching CSRF Token...");
    try {
      const response = await fetch('http://localhost:3000/api/auth/csrf', {
        credentials: 'include',
        method: 'GET',
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      // Type the response data
      const data: CsrfTokenResponse = await response.json();
      setCsrfToken(data.token);
      setApiResponse('CSRF Token fetched successfully.');
    } catch (error: any) { // Type the caught error
      console.error('Error fetching CSRF token:', error);
      setApiResponse('Error fetching token: ' + (error.message || 'Unknown error'));
    }
  };

  const callProtectedApi = async () => {
    console.log("Calling Protected API...");
    const token: string | null = getCsrfToken(); // Type the retrieved token

    if (!token) {
      console.warn("CSRF token not available. Cannot call protected API.");
      setApiResponse('Error: CSRF token not available. Click "Fetch and Set CSRF Token" first.');
      return;
    }

    try {
      const response = await fetch('/api/protected', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'X-CSRF-TOKEN': token,
        },
        body: JSON.stringify({ some: 'data' }),
      });

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }

      // Type the response data
      const data: ProtectedApiResponse = await response.json();
      setApiResponse('Protected API Response: ' + JSON.stringify(data, null, 2));
    } catch (error: any) { // Type the caught error
      console.error('Error calling protected API:', error);
      setApiResponse('Error calling protected API: ' + (error.message || 'Unknown error'));
    }
  };

  const handlePrintToken = () => {
    const currentToken: string | null = getCsrfToken(); // Type the retrieved token
    console.log("Current CSRF Token (from getCsrfToken):", currentToken);
  };

  const handleClearToken = () => {
      clearCsrfToken();
      setApiResponse('CSRF Token cleared.');
      console.log("CSRF Token cleared.");
  };

  return (
    <div>
      <h1>CSRF Token Test (TanStack Router Route - Type Safe)</h1>

      <button onClick={fetchCsrfToken}>Fetch and Set CSRF Token</button>
      <br></br>
      <button onClick={callProtectedApi}>Call Protected API</button>
      <br></br>
      <button onClick={handlePrintToken}>Print CSRF Token to Console</button>
      <br></br>
      <button onClick={handleClearToken}>Clear CSRF Token</button>

      <p>API Response: <pre>{apiResponse}</pre></p>
    </div>
  );
}
