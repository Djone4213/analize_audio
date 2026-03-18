import { useState } from "react";
import { useNavigate } from "react-router-dom";

export default function UploadFile() {
    const [file, setFile] = useState(null);
    const [thems, setThems] = useState([""]);
    const navigate = useNavigate();

    const handleAddThem = () => {
        setThems((prev) => [...prev, ""]);
    };

    const handleChangeThem = (index, value) => {
        const updated = [...thems];
        updated[index] = value;
        setThems(updated);
    };

    const handleSubmit = async (e) => {
        e.preventDefault();

        if (!file) return;

        const formData = new FormData();
        formData.append("video", file);

        thems.forEach((t) => {
            if (t.trim()) {
                formData.append("thems", t);
            }
        });

        await fetch("http://localhost:8080/api/audio", {
            method: "POST",
            body: formData,
        });

        // после загрузки возвращаемся на главную
        navigate("/");
        window.location.reload();
    };

    return (
        <div style={{ padding: 20 }}>
            <h1>Загрузка файла</h1>

            <button onClick={() => navigate("/")}>← Назад</button>

            <form onSubmit={handleSubmit} style={{ marginTop: 20 }}>
                <div>
                    <input
                        type="file"
                        accept="video/*"
                        onChange={(e) => setFile(e.target.files[0])}
                    />
                </div>

                <h3>Темы</h3>

                {thems.map((t, i) => (
                    <div key={i}>
                        <input
                            type="text"
                            value={t}
                            onChange={(e) => handleChangeThem(i, e.target.value)}
                            placeholder="Введите тему"
                        />
                    </div>
                ))}

                <button type="button" onClick={handleAddThem}>
                    + Добавить тему
                </button>

                <br /><br />

                <button type="submit">Отправить</button>
            </form>
        </div>
    );
}