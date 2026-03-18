import { useEffect, useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { getAudioById } from "../../services/audioApi";

export default function Details() {
    const { id } = useParams();
    const navigate = useNavigate();

    const [file, setFile] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState("");

    useEffect(() => {
        loadData();
    }, [id]);

    const loadData = async () => {
        try {
            setLoading(true);
            const data = await getAudioById(id);
            setFile(data);
        } catch (e) {
            setError(e.message);
        } finally {
            setLoading(false);
        }
    };

    return (
        <div style={{ padding: 20, maxWidth: 800, margin: "0 auto" }}>
            <button onClick={() => navigate(-1)}>
                ← Назад
            </button>

            {loading && <p>Загрузка...</p>}
            {error && <p style={{ color: "red" }}>{error}</p>}

            {file && (
                <>
                    <h2 style={{ marginTop: 20 }}>
                        {file.src_file_name}
                    </h2>

                    <div style={{ marginTop: 20 }}>
                        <label><b>Информация о сервисах:</b></label>
                        <textarea
                            value={file.analysis_text || ""}
                            readOnly
                            rows={10}
                            style={{
                                width: "100%",
                                marginTop: 5,
                                resize: "vertical"
                            }}
                        />
                    </div>

                    <div style={{ marginTop: 20 }}>
                        <label><b>Транскрибация звонка:</b></label>
                        <textarea
                            value={file.transcribed_text || ""}
                            readOnly
                            rows={10}
                            style={{
                                width: "100%",
                                marginTop: 5,
                                resize: "vertical"
                            }}
                        />
                    </div>
                </>
            )}
        </div>
    );
}