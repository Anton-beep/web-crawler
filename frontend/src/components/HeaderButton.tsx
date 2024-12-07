import {Link, useLocation} from "react-router-dom";
import {Button} from "@/components/ui/button.tsx";

export function HeaderButton({toLink, buttonString}: { toLink: string, buttonString: string }) {
    const location = useLocation();

    return (
        <>
            <Link to={toLink}>
                <Button variant="jumpingLink"
                        className={location.pathname === toLink ? "text-accent" : ""}>{buttonString}</Button>
            </Link>
        </>
    );
}