import { AppShell, Button, Stack } from "@mantine/core";
import NavButton from "./NavButton";
import { useAuth } from "../auth/AuthContext";

export interface SideNavProps {
    close: () => void;
}

export default function SideNav({ close }: SideNavProps) {
    const { loggedIn, signIn } = useAuth();
    return (
        <AppShell.Navbar>
            <Stack gap="md">
                {loggedIn ? (
                    <>
                        <NavButton to="/dashboard" onClick={close}>
                            Dashboard
                        </NavButton>
                        <NavButton to="/tax" onClick={close}>
                            Tax
                        </NavButton>
                        <NavButton to="/account" onClick={close}>
                            Account
                        </NavButton>
                    </>
                ) : (
                    <Button
                        onClick={() => {
                            signIn();
                            close();
                        }}
                        color="indigo"
                    >
                        Sign in
                    </Button>
                )}
            </Stack>
        </AppShell.Navbar>
    );
}
