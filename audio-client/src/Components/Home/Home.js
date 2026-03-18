import { useEffect, useState } from "react";
import { useNavigate } from "react-router-dom";
import { getAudioList } from "../../services/audioApi";

export default function Home() {
    const [files, setFiles] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");
    const navigate = useNavigate();

    useEffect(() => {
        loadFiles();
    }, []);

    const loadFiles = async () => {
        try {
            setLoading(true);
            const data = await getAudioList();
            setFiles(data);
        } catch (e) {
            setError(e.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{ padding: 20 }}>
            <h1>Загруженные файлы</h1>

            <button onClick={() => navigate("/upload")}>
                Загрузить файл
            </button>

            {loading && <p>Загрузка...</p>}
            {error && <p style={{ color: "red" }}>{error}</p>}

            <ul>
                {files.map((f) => (
                    <li key={f.id}>
                        {f.src_file_name}
                        {f.thems?.length > 0 && ` (${f.thems.join(", ")})`}
                    </li>
                ))}
            </ul>
        </div>
    );
}