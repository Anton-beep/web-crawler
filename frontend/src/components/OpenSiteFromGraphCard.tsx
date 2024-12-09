import {Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle} from "@/components/ui/dialog.tsx";
import {Button} from "@/components/ui/button.tsx";

const openSiteInNewWindow = (url: string) => {
    if (url) {
        window.open(url, '_blank');
    }
}

export default function OpenSiteFromGraphCard({url, setUrl}: { url: string, setUrl: (url: string) => void }) {
    return (
        <Dialog open={url !== ""} onOpenChange={() => {
            setUrl("")
        }}>
            <div className="mx-4">
                <DialogContent className="text-primary" onCloseAutoFocus={() => {
                    setUrl("");
                }}>
                    <DialogHeader>
                        <DialogTitle>Open URL</DialogTitle>
                        <DialogDescription>
                            Are you sure you want to open this URL: <span className="text-accent"> {url} </span> ?
                        </DialogDescription>
                    </DialogHeader>
                    <div className="grid gap-4 py-4">
                        <div className="grid grid-cols-4 items-center gap-4">
                            <Button className="col-span-4 bg-blue-600 text-primary" onClick={() => {
                                openSiteInNewWindow(url);
                            }}>
                                Open URL
                            </Button>
                        </div>
                    </div>
                </DialogContent>
            </div>
        </Dialog>
    );
}