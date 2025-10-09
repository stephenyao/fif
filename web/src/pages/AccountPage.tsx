import { Button, Group } from "@mantine/core";
import { useEffect } from "react";
import { useAuth } from "../auth/AuthContext";

export default function AccountPage() {
  const { user, signOut } = useAuth();

  // Fetch from backend when the page loads
  useEffect(() => {
    const controller = new AbortController();
    (async () => {
      if (!user) return; // require logged-in user
      try {
        const token = await user.getIdToken();
        const res = await fetch("http://localhost:8080/account", {
          method: "GET",
          headers: {
            Accept: "application/json",
            Authorization: `Bearer ${token}`,
          },
          signal: controller.signal,
        });
        const contentType = res.headers.get("content-type") || "";
        const body = contentType.includes("application/json")
          ? await res.json()
          : await res.text();
        console.log("GET /account response:", body);
      } catch (err: any) {
        if (err?.name !== "AbortError") {
          console.error("Failed to fetch /account:", err);
        }
      }
    })();
    return () => controller.abort();
  }, [user]);

  return (
    <>
      <h1>Account</h1>
      <p>Manage your account here.</p>
      <Group mt="md">
        <Button color="red" variant="filled" onClick={signOut}>
          Log out
        </Button>
      </Group>
    </>
  );
}
