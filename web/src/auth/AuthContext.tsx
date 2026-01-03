import {
    createContext,
    useContext,
    useEffect,
    useMemo,
    useState,
    type ReactNode,
} from "react";
import { supabase } from "../lib/supabase";
import type { User } from "@supabase/supabase-js";

interface AuthContextValue {
    loggedIn: boolean;
    user: User | null;
    initializing: boolean;
    signIn: () => Promise<void>;
    signOut: () => Promise<void>;
}

const AuthContext = createContext<AuthContextValue | undefined>(undefined);

export function AuthProvider({ children }: { children: ReactNode }) {
    const [user, setUser] = useState<User | null>(null);
    const [initializing, setInitializing] = useState(true);

    useEffect(() => {
        let mounted = true;

        supabase.auth
            .getSession()
            .then(({ data, error }) => {
                if (!mounted) return;
                if (error) {
                    console.error("Failed to get Supabase session:", error);
                }
                setUser(data.session?.user ?? null);
                setInitializing(false);
            })
            .catch((err) => {
                if (!mounted) return;
                console.error("Failed to initialize Supabase auth:", err);
                setInitializing(false);
            });

        const { data } = supabase.auth.onAuthStateChange((_event, session) => {
            setUser(session?.user ?? null);
            setInitializing(false);
        });

        return () => {
            mounted = false;
            data.subscription.unsubscribe();
        };
    }, []);

    const value = useMemo<AuthContextValue>(
        () => ({
            loggedIn: Boolean(user),
            user,
            initializing,
            signIn: async () => {
                const redirectTo =
                    import.meta.env.VITE_SUPABASE_REDIRECT_URL;
                if (!redirectTo) {
                    console.error(
                        "VITE_SUPABASE_REDIRECT_URL is required for sign-in."
                    );
                    return;
                }
                const { error } = await supabase.auth.signInWithOAuth({
                    provider: "google",
                    options: { redirectTo },
                });
                if (error) {
                    console.error("Supabase sign-in failed:", error);
                }
            },
            signOut: async () => {
                const { error } = await supabase.auth.signOut();
                if (error) {
                    console.error("Supabase sign-out failed:", error);
                }
            },
        }),
        [user, initializing]
    );

    // Optionally, render children immediately; Mantine UI can handle a momentary flash
    return (
        <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
    );
}

export function useAuth() {
    const ctx = useContext(AuthContext);
    if (!ctx) throw new Error("useAuth must be used within an AuthProvider");
    return ctx;
}
