export async function getAudioList() {
    const res = await fetch("http://localhost:8080/api/audio");

    if (!res.ok) {
        throw new Error("Ошибка загрузки файлов");
    }

    return res.json();
}