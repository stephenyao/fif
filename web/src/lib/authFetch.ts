import { supabase } from "./supabase";

export async function authFetch(
    input: RequestInfo | URL,
    init: RequestInit = {}
) {
    const { data, error } = await supabase.auth.getSession();
    if (error) {
        console.error("Failed to fetch Supabase session:", error);
    }
    const token = data.session?.access_token;
    const headers = new Headers(init.headers || {});
    if (token) headers.set("Authorization", `Bearer ${token}`);
    return fetch(input, { ...init, headers });
}
