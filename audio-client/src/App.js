import { BrowserRouter, Routes, Route } from "react-router-dom";
import Home from "./Components/Home/Home";
import UploadFile from "./Components/UploadFile/UploadFile";
import Details from "./Components/Details/Details";

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/upload" element={<UploadFile />} />
                <Route path="/details/:id" element={<Detailsпш />} />
            </Routes>
        </BrowserRouter>
    );
}

export default App;