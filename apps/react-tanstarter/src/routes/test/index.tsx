import { createFileRoute, redirect } from "@tanstack/react-router";

export const Route = createFileRoute("/test/")({
  component: RouteComponent,
  beforeLoad: async ({ context }) => {
    if (!context.user) {
      throw redirect({ to: "/signin" });
    }
  },
});

function RouteComponent() {
  return <div>Hello "/test/"!</div>;
}
