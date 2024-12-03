import {Button} from "@/components/ui/button.tsx";


export default function Header() {
    return (
        <header>
            <div>
                <h1 className="text-primary">Header</h1>
                <Button variant="link">Home</Button>
            </div>
        </header>
    )
}