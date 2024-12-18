import {Label} from "@/components/ui/label.tsx";
import {Input} from "@/components/ui/input.tsx";
import Api from "@/services/Api.ts";
import {useEffect, useState} from "react";
import {Button} from "@/components/ui/button.tsx";
import {isPasswordValid} from "@/utils/isPasswordValid.ts";
import {isUsernameValid} from "@/utils/isUsernameValid.ts";
import {isEmailValid} from "@/utils/isEmailValid.ts";
import Cookies from "js-cookie";
import {useNavigate} from "react-router-dom";

export default function Profile() {
    const [currentUsername, setCurrentUsername] = useState("");
    const [currentEmail, setCurrentEmail] = useState("");
    const [username, setUsername] = useState("");
    const [email, setEmail] = useState("");
    const [newPassword, setNewPassword] = useState("");
    const [currentPassword, setCurrentPassword] = useState("");
    const [message, setMessage] = useState("");
    const [isError, setIsError] = useState(false);
    const navigate = useNavigate();

    useEffect(() => {
        handleUserData();
    }, [currentUsername, currentEmail]);

    const handleUserData = () => {
        Api.getUser().then(response => {
            setCurrentUsername(response.data?.username);
            setCurrentEmail(response.data?.email);
            setUsername(response.data?.username);
            setEmail(response.data?.email);
        }).catch(err => console.log(err));
    }

    const handleSendData = () => {
        if (currentEmail === email && currentUsername === username && newPassword === "") {
            setMessage("You didn't change anything");
            return;
        }

        if (currentPassword === "") {
            setIsError(true);
            setMessage("You need to send your password if you want to change your data");
            return;
        }

        if (newPassword !== "" && !isPasswordValid(newPassword)) {
            setIsError(true);
            setMessage("Password must be at least 8 characters long and contain at least one uppercase letter, one lowercase letter, and one number.");
            return
        }

        if (!isUsernameValid(username)) {
            setIsError(true);
            setMessage("Username must be at least 4 characters long.");
            return
        }

        if (!isEmailValid(email)) {
            setIsError(true);
            setMessage("Email is not valid");
            return
        }

        Api.updateUser(username, email, newPassword, currentPassword).then(response => {
            if (response.status !== 200) {
                setIsError(true);
                setMessage("Error updating user data");
                console.error(response);
                return;
            }
            setIsError(false);
            setMessage("User data has been updated successfully");
            Cookies.set("access", response.data.access, { sameSite: 'Strict' });
            navigate(0);
        }).catch(err => {
            setIsError(true);
            setMessage("Error updating user data");
            console.log(err)
        });
    }

    const getMessage = () => {
        if (message !== "") {
            return (
                <div className={isError ? "text-red-500" : "text-green-500"}>
                    {message}
                </div>
            )
        }
    }

    return (
        <div className="flex justify-center items-center min-h-screen text-primary">
            <div className="grid gap-4 py-4">
                <div className="grid grid-cols-4 items-center gap-4">
                    <Label htmlFor="name" className="text-right">
                        Username
                    </Label>
                    <Input
                        type="text"
                        id="username"
                        className="col-span-3"
                        placeholder="write new username"
                        value={username === undefined ? "" : username}
                        onChange={e => setUsername(e.target.value)}
                    />
                    <div className="text-secondary">Current username: {currentUsername}</div>
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                    <Label htmlFor="url" className="text-right">
                        Email
                    </Label>
                    <Input
                        type="email"
                        id="email"
                        placeholder="write new email"
                        className="col-span-3"
                        value={email === undefined ? "" : email}
                        onChange={e => setEmail(e.target.value)}
                    />
                    <div className="text-secondary">Current email: {currentEmail}</div>
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                <Label htmlFor="url" className="text-right">
                        New Password
                    </Label>
                    <Input
                        type="password"
                        id="newPassword"
                        placeholder="write new password"
                        className="col-span-3"
                        onChange={e => setNewPassword(e.target.value)}
                    />
                </div>
                <br className="mb-4"></br>
                <div className="text-accent">You need to send your password if you want to change your data</div>
                <div className="grid grid-cols-4 items-center gap-4">
                    <Label htmlFor="url" className="text-right">
                        Current Password
                    </Label>
                    <Input
                        type="password"
                        id="currentPassword"
                        placeholder="write current password"
                        className="col-span-3"
                        onChange={e => setCurrentPassword(e.target.value)}
                    />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                    <Button variant="jumping" className="bg-green-600 text-primary mx-4"
                            onClick={handleSendData}>Update</Button>
                </div>
                <br className="mb-4"></br>
                {getMessage()}
            </div>
        </div>
    )
}