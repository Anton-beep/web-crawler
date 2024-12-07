import {HeaderButton} from "@/components/HeaderButton.tsx";


export default function Header() {
    return (
        <header className="mb-8 mx-4">
            <div>
                <h1 className="text-primary font-bold text-2xl mb-1">Web Crawler (we need a logo)</h1>
                <HeaderButton toLink="/" buttonString="About"/>
                <HeaderButton toLink="/dashboard" buttonString="Dashboard"/>
            </div>
        </header>
    )
}