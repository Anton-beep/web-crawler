import {DialogContent, DialogDescription, DialogFooter, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Label} from "@/components/ui/label.tsx";
import {Input} from "@/components/ui/input.tsx";
import {Button} from "@/components/ui/button.tsx";
import {useState} from "react";
import Api from "@/services/Api.ts";
import checkValidUrl from "@/utils/checkValidUrl.ts";


export default function CreateProjectCard() {
    const [name, setName] = useState("");
    const [url, setUrl] = useState("");
    const [message, setMessage] = useState("");
    const [isError, setIsError] = useState(false);

    const handleSend = () => {

        if (!checkValidUrl(url)) {
            setIsError(true);
            setMessage("Invalid url");
            return;
        }

        Api.createProject(name, url).then(() => {
            setIsError(false);
            setMessage("Project created successfully");
            setName("");
            setUrl("");
        }).catch((e) => {
            setIsError(true);
            setMessage("Error creating project");
            console.error(e);
        })
    }

    return (
        <DialogContent className="text-primary">
            <DialogHeader>
                <DialogTitle>Create Project</DialogTitle>
                <DialogDescription>
                    Give a link to the website you want to crawl
                </DialogDescription>
            </DialogHeader>
            <div className="grid gap-4 py-4">
                <div className="grid grid-cols-4 items-center gap-4">
                    <Label htmlFor="name" className="text-right">
                        Name
                    </Label>
                    <Input
                        type="text"
                        id="name"
                        className="col-span-3"
                        placeholder="Name your project"
                        value={name}
                        onChange={e => setName(e.target.value)}
                    />
                </div>
                <div className="grid grid-cols-4 items-center gap-4">
                    <Label htmlFor="url" className="text-right">
                        Start Url
                    </Label>
                    <Input
                        type="url"
                        id="url"
                        placeholder="start url"
                        className="col-span-3"
                        value={url}
                        onChange={e => setUrl(e.target.value)}
                    />
                </div>
            </div>
            <DialogFooter>
                <div>
                    {isError ? <p className="text-red-500">{message}</p> : message &&
                        <p className="text-green-500">{message}</p>}
                </div>
                <div>
                    <Button variant="jumping" className="bg-green-600 text-primary mx-4" onClick={handleSend}>Create
                        Project</Button>
                </div>
            </DialogFooter>
        </DialogContent>
    )
}