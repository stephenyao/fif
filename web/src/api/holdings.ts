import { authFetch } from "../lib/authFetch";
import type { Holding } from "../models/Holding";

export async function getHoldings(signal?: AbortSignal): Promise<Holding[]> {
    const res = await authFetch("http://localhost:8080/holdings", {
        method: "GET",
        signal,
    });
    if (!res.ok) throw new Error(`HTTP ${res.status}`);

    const data = await res.json();
    return data as Holding[];
}
