import {DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Label} from "@/components/ui/label.tsx";
import {Input} from "@/components/ui/input.tsx";
import {Button} from "@/components/ui/button.tsx";
import {useState} from "react";
import Api from "@/services/Api.ts";
import Cookies from "js-cookie";
import {useNavigate} from "react-router-dom";


export default function LoginCard({setIsOpen}: { setIsOpen: (isOpen: boolean) => void }) {
    const [login, setLogin] = useState("");
    const [password, setPassword] = useState("");
    const [message, setMessage] = useState("");
    const [isError, setIsError] = useState(false);
    const navigate = useNavigate();

    const handleSend = () => {
        Api.loginUser(login, password).then((response) => {
            if (response.status !== 200) {
                setIsError(true);
                setMessage("Error logging in");
                console.error(response);
                return;
            }
            Cookies.set("access", response.data.access, { sameSite: 'Strict' });
            console.log(Cookies.get());
            setIsError(false);
            setMessage("You have been logged in successfully");
            setIsOpen(false);
            setLogin("");
            setPassword("");
            navigate(0);
        }).catch((e) => {
            setIsError(true);
            setMessage("Invalid login or password");
            console.error(e);
        })
    }

    return (
        <DialogContent className="text-primary">
            <DialogHeader>
                <DialogTitle>Login</DialogTitle>
                <DialogDescription>
                    Please fill in the form to login
                </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
                <div className="grid grid-cols-4 items-center gap-4">
                    <Label htmlFor="name" className="text-right">
                        Login
                    </Label>
                    <Input
                        type="text"
                        id="login"
                        className="col-span-3"
                        placeholder="your username or email"
                        value={login}
                        onChange={e => setLogin(e.target.value)}
                    />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                    <Label htmlFor="url" className="text-right">
                        Password
                    </Label>
                    <Input
                        type="password"
                        id="password"
                        placeholder="your password"
                        className="col-span-3"
                        value={password}
                        onChange={e => setPassword(e.target.value)}
                    />
                </div>
            </div>
            <DialogFooter>
                <div>
                    {isError ? <p className="text-red-500">{message}</p> : message &&
                        <p className="text-green-500">{message}</p>}
                </div>
                <div>
                    <Button variant="jumping" className="bg-green-600 text-primary mx-4" onClick={handleSend}>Login</Button>
                </div>
            </DialogFooter>
        </DialogContent>
    )
}