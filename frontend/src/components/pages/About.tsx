import {MacbookScroll} from "@/components/ui/macbook-scroll.tsx";
import AnimationOnScroll from "@/components/AnimationOnScroll.tsx";

export default function About() {
    return (
        <div>
            <MacbookScroll src="about_mackbook_content.png"
                           title={(<div className="text-primary">We crawl around the web to collect data for
                               <span className="text-accent font-extrabold"> you </span></div>)}/>

            <div className="flex flex-col md:flex-row items-center gap-8 mb-52 text-primary">
                <div className="md:w-1/2">
                    <img alt={"we need a picture here"} src="about_mackbook_content.png" className="rounded border border-zinc-600"/>
                </div>
                <div className="md:w-1/2">
                    <AnimationOnScroll startState={"transform -translate-x-3/4"} endState={"transform translate-x-0"}>
                        <p className="text-center text-3xl text-primary font-bold mx-32">We aim to create a robust tool
                            that
                            simplifies and enhances the <span className="text-accent">data collection process, empowering businesses and analysts to make
                        informed decisions.</span></p>
                    </AnimationOnScroll>
                </div>
            </div>

            <div className="flex flex-col md:flex-row items-center gap-8 mb-52 text-primary">
                <div className="md:w-1/2">
                    <AnimationOnScroll startState={"transform translate-x-3/4"} endState={"transform translate-x-0"}>
                        <p className="text-center text-3xl text-primary font-bold mx-32">Our crawler operates on a
                            microservice architecture, ensuring <span className="text-accent">high performance through horizontal scalability.</span> Each
                            page is processed by dedicated services for optimal efficiency.
                        </p>
                    </AnimationOnScroll>
                </div>
                <div className="md:w-1/2">
                    <img alt={"we need a picture here"} src="about_mackbook_content.png" className="rounded border border-zinc-600"/>
                </div>
            </div>
        </div>
    )
}