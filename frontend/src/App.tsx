import {BrowserRouter as Router, Route, Routes} from 'react-router-dom';
import Header from "./components/Header.tsx";
import Footer from "./components/Footer.tsx";
import About from "./components/pages/About.tsx";
import Dashboard from "./components/pages/Dashboard.tsx";
import Project from "@/components/pages/Project.tsx";
import Profile from "@/components/pages/Profile.tsx";


function App() {

    return (
        <Router>
            <div className="flex flex-col min-h-screen">
                <Header/>
                <div className="flex-auto mx-4">
                    <Routes>
                        <Route path="/" element={<About/>}/>
                        <Route path="/dashboard" element={<Dashboard/>}/>
                        <Route path="project/:projectId" element={<Project/>}/>
                        <Route path="/profile" element={<Profile/>}/>
                        <Route path="*" element={<About/>}/>
                    </Routes>
                </div>
                <Footer/>
            </div>
        </Router>
    );
}

export default App;