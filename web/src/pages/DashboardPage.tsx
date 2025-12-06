import {
    Card,
    Divider,
    Group,
    Paper,
    Skeleton,
    Stack,
    Table,
    Text,
    useMantineTheme,
} from "@mantine/core";
import { useMediaQuery } from "@mantine/hooks";
import { useQuery } from "@tanstack/react-query";
import { authFetch } from "../lib/authFetch";

type Holding = {
    name: string;
    symbol: string;
    quantity: number;
    currency: string;
    cost: string;
};

const fetchHoldings = async (): Promise<Holding[]> => {
    const res = await authFetch("http://localhost:8080/holdings", {
        method: "GET",
    });
    if (!res.ok) {
        throw new Error("Failed to fetch holdings");
    }

    return res.json();
};

export default function DashboardPage() {
    const {
        data: holdings,
        isLoading,
        isError,
        error,
    } = useQuery<Holding[], Error>({
        queryKey: ["holdings"],
        queryFn: fetchHoldings,
    });

    return (
        <>
            <h1>Dashboard</h1>
            <p>Welcome back! This is your dashboard.</p>

            {/* Loading */}
            {isLoading && (
                <Stack gap="md" mt="md" align="stretch">
                    {[1, 2, 3, 4].map((i) => (
                        <Stack gap={4} key={i}>
                            <Skeleton height={14} w="40%" />
                            <Skeleton height={18} w="100%" />
                        </Stack>
                    ))}
                </Stack>
            )}

            {/* Error */}
            {isError && (
                <div style={{ color: "red", marginTop: "1rem" }}>
                    Error! {error.message}
                </div>
            )}

            {/* Success */}
            {!isLoading && !isError && holdings && <Demo holdings={holdings} />}
        </>
    );
}

function Demo({ holdings }: { holdings: Holding[] }) {
    const theme = useMantineTheme();
    const isMobile = useMediaQuery(`(max-width: ${theme.breakpoints.sm})`);

    if (isMobile) {
        return <HoldingsCards holdings={holdings} />;
    }

    return <HoldingsTable holdings={holdings} />;
}

function HoldingsCards({ holdings }: { holdings: Holding[] }) {
    return (
        <Stack>
            {holdings.map((holding) => (
                <Card
                    shadow="sm"
                    p="md"
                    key={holding.symbol}
                    withBorder
                    radius="lg"
                >
                    <Group justify="space-between">
                        <Text fw={600}>{holding.name}</Text>
                        <Text fw={600} size="sm">
                            {holding.symbol}
                        </Text>
                    </Group>
                    <Divider />
                    <Group justify="space-between" mt="sm">
                        <Text size="sm" c="dimmed">
                            Currency
                        </Text>
                        <Text fw={500}>{holding.currency}</Text>
                    </Group>
                    <Group justify="space-between">
                        <Text size="sm" c="dimmed">
                            Quantity
                        </Text>
                        <Text fw={500}>{holding.quantity}</Text>
                    </Group>
                    <Group justify="space-between">
                        <Text size="sm" c="dimmed">
                            Cost (NZD)
                        </Text>
                        <Text fw={500}>{holding.cost}</Text>
                    </Group>
                </Card>
            ))}
        </Stack>
    );
}

function HoldingsTable({ holdings }: { holdings: Holding[] }) {
    return (
        <Paper shadow="sm" p="xl" radius="md">
            <Table>
                <Table.Thead>
                    <Table.Tr>
                        <Table.Th>Holding</Table.Th>
                        <Table.Th>Quantity</Table.Th>
                        <Table.Th>Currency</Table.Th>
                        <Table.Th>Cost (NZD)</Table.Th>
                    </Table.Tr>
                </Table.Thead>
                <Table.Tbody>
                    {holdings.map((holding) => (
                        <Table.Tr key={holding.symbol}>
                            <Table.Td>
                                {holding.name}{" "}
                                <strong>({holding.symbol})</strong>
                            </Table.Td>
                            <Table.Td>{holding.quantity}</Table.Td>
                            <Table.Td>{holding.currency}</Table.Td>
                            <Table.Td>{holding.cost}</Table.Td>
                        </Table.Tr>
                    ))}
                </Table.Tbody>
            </Table>
        </Paper>
    );
}
