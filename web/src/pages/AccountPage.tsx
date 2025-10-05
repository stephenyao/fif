import { Button, Group } from "@mantine/core";
import { useAuth } from "../auth/AuthContext";

export default function AccountPage() {
  const { signOut } = useAuth();
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
