import {DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Label} from "@/components/ui/label.tsx";
import {Input} from "@/components/ui/input.tsx";
import {Button} from "@/components/ui/button.tsx";
import {useState} from "react";
import Api from "@/services/Api.ts";
import Cookies from "js-cookie";
import {useNavigate} from "react-router-dom";
import {isPasswordValid} from "@/utils/isPasswordValid.ts";
import {isUsernameValid} from "@/utils/isUsernameValid.ts";


export default function RegistrationCard({setIsOpen}: { setIsOpen: (isOpen: boolean) => void }) {
    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [message, setMessage] = useState("");
    const [isError, setIsError] = useState(false);
    const navigate = useNavigate();

    const setCredentialsEmpty = () => {
        setUsername("");
        setEmail("");
        setPassword("");
    }

    const handleSend = () => {
        if (!isPasswordValid(password)) {
            setIsError(true);
            setMessage("Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, one number and one special character");
            return;
        }

        if (!isUsernameValid(username)) {
            setIsError(true);
            setMessage("Username must be more than 3 characters long.");
            return;
        }

        Api.registerUser(username, email, password).then((response) => {
            if (response.status !== 200) {
                setIsError(true);
                setMessage("Error registering user");
                console.error(response);
                return;
            }
            Cookies.set("access", response.data.access, { sameSite: 'Strict' });
            setIsError(false);
            setMessage("You have been registered successfully");
            setCredentialsEmpty();
            setIsOpen(false);
            navigate(0);
        }).catch((e) => {
            setIsError(true);
            if (e.data?.mesage !== undefined) {
                setMessage(e.data?.mesage);
            } else {
                setMessage("Error registering user");
            }
            console.error(e);
        })
    }

    return (
        <DialogContent className="text-primary">
            <DialogHeader>
                <DialogTitle>Register</DialogTitle>
                <DialogDescription>
                    Please fill in the form to register
                </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
                <div className="grid grid-cols-4 items-center gap-4">
                    <Label htmlFor="name" className="text-right">
                        Username
                    </Label>
                    <Input
                        type="text"
                        id="username"
                        className="col-span-3"
                        placeholder="your username"
                        value={username}
                        onChange={e => setUsername(e.target.value)}
                    />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                    <Label htmlFor="url" className="text-right">
                        Email
                    </Label>
                    <Input
                        type="email"
                        id="email"
                        placeholder="your email"
                        className="col-span-3"
                        value={email}
                        onChange={e => setEmail(e.target.value)}
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
                    <Button variant="jumping" className="bg-green-600 text-primary mx-4" onClick={handleSend}>Register</Button>
                </div>
            </DialogFooter>
        </DialogContent>
    )
}