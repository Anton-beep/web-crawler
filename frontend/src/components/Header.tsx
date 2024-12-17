import {HeaderButton} from "@/components/HeaderButton.tsx";
import {useEffect, useState} from "react";
import Api from "@/services/Api.ts";
import {Button} from "@/components/ui/button.tsx";
import RegistrationCard from "@/components/RegistrationCard.tsx";
import LoginCard from "@/components/LoginCard.tsx";
import {Dialog, DialogTrigger} from "@/components/ui/dialog.tsx";
import LogoutCard from "@/components/LogoutCard.tsx";

export default function Header() {
    const [user, setUser] = useState<string | undefined>(undefined);
    const [isRegistrationOpen, setIsRegistrationOpen] = useState(false);
    const [isLoginOpen, setIsLoginOpen] = useState(false);

    useEffect(() => {
        Api.getUser().then(response =>
            setUser(response.data?.username)
        ).catch(err => err.status !== 401 && console.log(err));
    }, [user]);

    const getUserButtons = () => {
        if (user !== undefined) {
            return (
                <div className="flex items-center">
                    <HeaderButton toLink="/profile" buttonString="Edit Profile"/>
                    <Dialog>
                        <DialogTrigger>
                            <Button variant="jumping" className="text-black mx-4">Log out</Button>
                        </DialogTrigger>
                        <LogoutCard/>
                    </Dialog>
                </div>
            )
        } else {
            return (
                <div className="flex items-center">
                    <Dialog open={isRegistrationOpen} onOpenChange={setIsRegistrationOpen}>
                        <DialogTrigger>
                            <Button variant="jumping" className="text-black mx-4"
                                    onClick={() => setIsRegistrationOpen(true)}>Register</Button>
                        </DialogTrigger>
                        <RegistrationCard setIsOpen={setIsRegistrationOpen}/>
                    </Dialog>

                    <Dialog open={isLoginOpen} onOpenChange={setIsLoginOpen}>
                        <DialogTrigger>
                            <Button variant="jumping" className="text-black mx-4">Log in</Button>
                        </DialogTrigger>
                        <LoginCard setIsOpen={setIsLoginOpen}/>
                    </Dialog>
                </div>
            )
        }
    }

    const getHeaderButtons = () => {
        if (user !== undefined) {
            return (
                <>
                    <HeaderButton toLink="/" buttonString="About"/>
                    <HeaderButton toLink="/dashboard" buttonString="Dashboard"/>
                </>
            )
        }

        return (
            <>
                <HeaderButton toLink="/" buttonString="About"/>
            </>
        )
    }

    return (
        <header className="mb-8 mx-4">
            <div className="flex items-center justify-between">
                <div className="flex items-center">
                    <img src="/web_crawler_logo.svg" alt="logo" className="w-20 h-20 mr-2 mt-2"
                         style={{filter: "invert(100%) sepia(0%) saturate(0%) hue-rotate(0deg) brightness(100%) contrast(100%)"}}/>
                    <h1 className="text-primary font-bold text-2xl mb-1">Web Crawler</h1>
                </div>

                <div className="flex items-center">
                    {getHeaderButtons()}
                </div>
                <div className="flex items-center">
                    {getUserButtons()}
                </div>
            </div>
        </header>
    )
}