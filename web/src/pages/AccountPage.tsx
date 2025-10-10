import { Alert, Button, Group, Paper, Stack, Text, Title } from "@mantine/core";
import { useEffect, useState } from "react";
import { useAuth } from "../auth/AuthContext";
import { getAccountProfile } from "../api/account";
import type { AccountProfile } from "../models/Account";

export default function AccountPage() {
    const { user, signOut } = useAuth();
    const [profile, setProfile] = useState<AccountProfile | null>(null);
    const [error, setError] = useState<string | null>(null);

    // Fetch from backend when the page loads
    useEffect(() => {
        const controller = new AbortController();
        (async () => {
            if (!user) return; // require logged-in user
            try {
                const profile = await getAccountProfile(controller.signal);
                setProfile(profile);
            } catch (err: any) {
                if (err?.name !== "AbortError") {
                    console.error("Failed to fetch /account:", err);
                    setError(err?.message || "Failed to load account");
                }
            }
        })();
        return () => controller.abort();
    }, [user]);

    return (
        <>
            <h1>Account</h1>

            {profile && (
                <Stack gap="sm">
                    <Stack gap={2}>
                        <Text size="sm" fw={600} c="dimmed">
                            Name
                        </Text>
                        <Text>{profile.name}</Text>
                    </Stack>

                    <Stack gap={2}>
                        <Text size="sm" fw={600} c="dimmed">
                            Email
                        </Text>
                        <Text>{profile.email}</Text>
                    </Stack>

                    {/* Add more attributes as needed */}
                </Stack>
            )}
            {error && (
                <Alert color="red" mt="md">
                    <Text c="red">{error}</Text>
                </Alert>
            )}
            <Group mt="md">
                <Button color="red" variant="filled" onClick={signOut}>
                    Log out
                </Button>
            </Group>
        </>
    );
}
