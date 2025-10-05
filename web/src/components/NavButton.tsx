import { Button } from "@mantine/core";
import { Link, matchPath, useLocation } from "react-router-dom";
import type { ReactNode } from "react";

function useIsActive(to: string) {
  const location = useLocation();
  const isRoot = to === "/";
  const active = isRoot
    ? location.pathname === "/"
    : Boolean(
        matchPath({ path: to, caseSensitive: false, end: false }, location.pathname)
      );
  return active;
}

export interface NavButtonProps {
  to: string;
  children: ReactNode;
  visibleFrom?: string;
  hiddenFrom?: string;
  onClick?: () => void;
}

export default function NavButton({
  to,
  children,
  visibleFrom,
  hiddenFrom,
  onClick,
}: NavButtonProps) {
  const active = useIsActive(to);
  return (
    <Button
      component={Link}
      to={to}
      onClick={onClick}
      variant={active ? "filled" : "subtle"}
      color="indigo"
      aria-current={active ? "page" : undefined}
      visibleFrom={visibleFrom}
      hiddenFrom={hiddenFrom}
    >
      {children}
    </Button>
  );
}

