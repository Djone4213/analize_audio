import { BrowserRouter, Routes, Route } from "react-router-dom";
import Home from "./Components/Home/Home";
import UploadFile from "./Components/UploadFile/UploadFile";

function App() {
    return (
        <BrowserRouter>
            <Routes>
                <Route path="/" element={<Home />} />
                <Route path="/upload" element={<UploadFile />} />
            </Routes>
        </BrowserRouter>
    );
}

export default App;