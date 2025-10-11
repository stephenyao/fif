import { Alert, Button, Group, Skeleton, Stack, Text } from "@mantine/core";
import { useEffect, useState } from "react";
import { useAuth } from "../auth/AuthContext";
import { getAccountProfile } from "../api/account";
import type { AccountProfile } from "../models/Account";

export default function AccountPage() {
    const { user, signOut } = useAuth();
    const [profile, setProfile] = useState<AccountProfile | null>(null);
    const [error, setError] = useState<string | null>(null);
    const [loading, setLoading] = useState<boolean>(true);

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
            setLoading(false);
        })();
        return () => controller.abort();
    }, [user]);

    return (
        <>
            <h1>Account</h1>

            <AccountProfileSection
                loading={loading}
                profile={profile}
                error={error}
            />

            <Group mt="md">
                <Button color="red" variant="filled" onClick={signOut}>
                    Log out
                </Button>
            </Group>
        </>
    );
}

interface AccountProfileSectionProps {
    loading: boolean;
    profile: AccountProfile | null;
    error: string | null;
}

function AccountProfileSection({
    loading,
    profile,
    error,
}: AccountProfileSectionProps) {
    return loading ? (
        <Stack gap="md" mt="md" align="stretch" w={400}>
            {[1, 2].map((i) => (
                <Stack gap={4} key={i}>
                    <Skeleton height={14} w="40%" />
                    <Skeleton height={18} w="100%" />
                </Stack>
            ))}
        </Stack>
    ) : (
        <>
            {profile && (
                <Stack gap="sm" mt="md">
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
                </Stack>
            )}

            {error && (
                <Alert color="red" mt="md">
                    <Text c="red">{error}</Text>
                </Alert>
            )}
        </>
    );
}
