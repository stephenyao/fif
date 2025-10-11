import { AppShell, Burger, Button, Group, Text } from "@mantine/core";
import NavButton from "./NavButton";
import { useAuth } from "../auth/AuthContext";

export interface TopNavProps {
    opened: boolean;
    toggle: () => void;
}

export default function TopNav({ opened, toggle }: TopNavProps) {
    const { loggedIn, signIn } = useAuth();
    return (
        <AppShell.Header>
            <Group h="100%" px="md" justify="space-between">
                <Text fw={700}>MyApp</Text>

                {loggedIn ? (
                    <Group gap="md" visibleFrom="md">
                        <NavButton to="/dashboard">Dashboard</NavButton>
                        <NavButton to="/tax">Tax</NavButton>
                    </Group>
                ) : (
                    <div />
                )}

                <Group>
                    <Burger
                        opened={opened}
                        onClick={toggle}
                        hiddenFrom="md"
                        size="sm"
                        aria-label={
                            opened ? "Close navigation" : "Open navigation"
                        }
                    />
                    {loggedIn ? (
                        <NavButton to="/account" visibleFrom="md">
                            Account
                        </NavButton>
                    ) : (
                        <Button
                            onClick={signIn}
                            visibleFrom="md"
                            color="indigo"
                        >
                            Sign in
                        </Button>
                    )}
                </Group>
            </Group>
        </AppShell.Header>
    );
}
