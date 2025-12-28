// App.tsx
// AppShell layout and routes (providers live in main.tsx)

import "@mantine/core/styles.css";
import { AppShell, Center, Container, Loader } from "@mantine/core";
import { Routes, Route, Navigate } from "react-router-dom";
import { useDisclosure } from "@mantine/hooks";
import TopNav from "./components/TopNav";
import SideNav from "./components/SideNav";
import HomePage from "./pages/HomePage";
import TaxPage from "./pages/TaxPage";
import AccountPage from "./pages/AccountPage";
import DashboardPage from "./pages/DashboardPage";
import ProtectedRoute from "./auth/ProtectedRoute";
import { useAuth } from "./auth/AuthContext";

export default function App() {
    const { loggedIn, initializing } = useAuth();
    const [opened, { toggle, close }] = useDisclosure();
    if (initializing) {
        return (
            <Center style={{ position: "fixed", inset: 0 }}>
                <Loader color="indigo" size="lg" />
            </Center>
        );
    }
    return (
        <AppShell
            padding="md"
            header={{ height: 60 }}
            navbar={{
                width: 300,
                breakpoint: "md",
                collapsed: { mobile: !opened, desktop: true },
            }}
        >
            <TopNav opened={opened} toggle={toggle} />
            <SideNav close={close} />
            <AppShell.Main>
                <Container size="lg">
                    <Routes>
                        <Route
                            path="/"
                            element={
                                loggedIn ? (
                                    <Navigate to="/dashboard" replace />
                                ) : (
                                    <HomePage />
                                )
                            }
                        />
                        <Route
                            path="/dashboard"
                            element={
                                <ProtectedRoute>
                                    <DashboardPage />
                                </ProtectedRoute>
                            }
                        />
                        <Route
                            path="/tax"
                            element={
                                <ProtectedRoute>
                                    <TaxPage />
                                </ProtectedRoute>
                            }
                        />
                        <Route
                            path="/account"
                            element={
                                <ProtectedRoute>
                                    <AccountPage />
                                </ProtectedRoute>
                            }
                        />
                        <Route path="*" element={<HomePage />} />
                    </Routes>
                </Container>
            </AppShell.Main>
        </AppShell>
    );
}
