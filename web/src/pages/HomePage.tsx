import { Button, Group, Stack, Text } from "@mantine/core";
import { useAuth } from "../auth/AuthContext";

export default function HomePage() {
  const { loggedIn, signIn } = useAuth();
  return (
    <>
      <h1>Home</h1>
      <p>Welcome to the home page.</p>
      {!loggedIn && (
        <Stack mt="lg" align="center">
          <Text>Sign in to access your dashboard and taxes.</Text>
          <Group>
            <Button color="indigo" onClick={signIn}>
              Sign in with Google
            </Button>
          </Group>
        </Stack>
      )}
    </>
  );
}
