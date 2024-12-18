import {Card, CardContent, CardFooter, CardHeader} from "@/components/ui/card.tsx";
import {Button} from "@/components/ui/button.tsx";
import Cookies from "js-cookie";

export function HelpCard({setShowHelp}: { setShowHelp: (show: boolean) => void }) {
    return (
        <Card className="mb-12 text-lg border-0">
            <CardHeader className="font-extrabold text-xl">
                Are you new here?
            </CardHeader>
            <CardContent>
                <p className="text-lg">Welcome to our app! Here are some tips to get for you to start:</p>
                <ul className="list-disc list-inside">
                    <li>Create a new project by clicking the <span className="text-accent">"Create Project"</span> button. You need to input a name for the project, start url, max amount of links in your project, and depth of the project. The depth of the project refers to how many levels deep the web crawler will go when following links from the start URL.</li>
                    <li>Open an existing project by clicking the <span className="text-accent">"Open Project"</span> button. You will see a list of your projects. Click on the button to see the details.</li>
                </ul>
            </CardContent>
            <CardFooter>
                <Button variant="outline" onClick={() => {
                    setShowHelp(false);
                    Cookies.set("showHelp", "false");
                }}>Got it!</Button>
            </CardFooter>
        </Card>
    )
}