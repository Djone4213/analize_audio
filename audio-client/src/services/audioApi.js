export async function getAudioList() {
    const res = await fetch("http://localhost:8080/api/audio");

    if (!res.ok) {
        throw new Error("Ошибка загрузки файлов");
    }

    return res.json();
}

export async function getAudioById(id) {
    const res = await fetch(`http://localhost:8080/api/audio/${id}`);

    if (!res.ok) {
        throw new Error("Ошибка загрузки файлов");
    }

    return res.json();
}