import {DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Button} from "@/components/ui/button.tsx";
import Cookies from "js-cookie";
import {useNavigate} from "react-router-dom";


export default function LogoutCard() {
    const navigate = useNavigate();

    const handleLogout = () => {
        Cookies.remove("access", {sameSite: 'Strict'});
        navigate(0);
    }

    return (
        <DialogContent className="text-primary">
            <DialogHeader>
                <DialogTitle>Log out</DialogTitle>
                <DialogDescription>
                    Are you sure you want to log out?
                </DialogDescription>
            </DialogHeader>
            <DialogFooter>
                <div>
                    <Button variant="jumping" className="bg-red-800 text-primary mx-4" onClick={handleLogout}>Log
                        out</Button>
                </div>
            </DialogFooter>
        </DialogContent>
    )
}