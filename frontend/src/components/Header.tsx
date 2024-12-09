import {HeaderButton} from "@/components/HeaderButton.tsx";


export default function Header() {
    return (
        <header className="mb-8 mx-4">
            <div>
                <div className="flex items-center">
                    <img src="/web_crawler_logo.svg" alt="logo" className="w-20 h-20 mr-2"
                         style={{filter: "invert(100%) sepia(0%) saturate(0%) hue-rotate(0deg) brightness(100%) contrast(100%)"}}/>
                    <h1 className="text-primary font-bold text-2xl mb-1">Web Crawler</h1>
                </div>

                <HeaderButton toLink="/" buttonString="About"/>
                <HeaderButton toLink="/dashboard" buttonString="Dashboard"/>
            </div>
        </header>
    )
}