import { authFetch } from "../lib/authFetch";
import { type AccountProfile } from "../models/Account";

export async function getAccountProfile(
    signal?: AbortSignal
): Promise<AccountProfile> {
    const res = await authFetch(`${import.meta.env.VITE_API_URL}/account`, {
        method: "GET",
        signal,
    });
    if (!res.ok) throw new Error(`HTTP ${res.status}`);

    const data = await res.json();
    return data as AccountProfile;
}
