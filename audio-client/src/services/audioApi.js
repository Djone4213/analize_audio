const API_URL = import.meta.env.VITE_API_URL;

export async function getAudioList() {
    const res = await fetch(`${API_URL}/api/audio`);

    if (!res.ok) {
        throw new Error("Ошибка загрузки файлов");
    }

    return res.json();
}

export async function getAudioById(id) {
    const res = await fetch(`${API_URL}/api/audio/${id}`);

    if (!res.ok) {
        throw new Error("Ошибка загрузки файлов");
    }

    return res.json();
}