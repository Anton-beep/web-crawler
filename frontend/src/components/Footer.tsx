import {Button} from "@/components/ui/button.tsx";
import {Dialog, DialogContent, DialogHeader, DialogTitle, DialogTrigger} from "@/components/ui/dialog.tsx";

export default function Footer() {
    return (
        <footer className="w-full">
            <div className="bg-gray-800 text-primary py-4">
                <div className="container flex justify-between w-full">
                    <div className="text-left ">
                        <span className="ml-10">Web Crawler v1.0.0</span>
                        <Dialog>
                            <DialogTrigger>
                                <Button variant="jumpingLink">
                                    Creators
                                </Button>
                            </DialogTrigger>
                            <DialogContent>
                                <DialogHeader className="text-primary font-extrabold">
                                    <DialogTitle>Creators</DialogTitle>
                                </DialogHeader>
                                <div className="text-accent underline">
                                    <div className="my-5">
                                        <a href="https://github.com/Cyber-Zhaba">Cyber-Zhaba</a>
                                    </div>
                                    <div className="my-5">
                                        <a href="https://github.com/Anton-beep">Anton-beep</a>
                                    </div>
                                    <div className="my-5">
                                        <a href="https://github.com/DrPepper1337">DrPepper1337</a>
                                    </div>
                                </div>
                            </DialogContent>
                        </Dialog>
                    </div>
                </div>
            </div>
        </footer>
    )
}