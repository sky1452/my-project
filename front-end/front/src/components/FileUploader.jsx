import { useState, useEffect } from "react";
import { useDropzone } from "react-dropzone";
import {
  Trash2, FileText
} from "lucide-react";

export default function FileUploader({ onFilesChange }) {
  const [files, setFiles] = useState([]);

  const onDrop = (acceptedFiles) => {
    const mapped = acceptedFiles.map(file => ({
      id: Math.random().toString(36).slice(2),
      file,
      progress: 0,
      status: null
    }));

    setFiles(prev => [...prev, ...mapped]);
  };

  const removeFile = (id) => {
    setFiles(prev => prev.filter(f => f.id !== id));
  };

  // ✅ ВАЖНО: теперь обновляем родителя здесь
  useEffect(() => {
    onFilesChange(files.map(f => f.file));
  }, [files, onFilesChange]);

  const { getRootProps, getInputProps } = useDropzone({
    onDrop,
    accept: {
      "application/pdf": [".pdf"]
    },
    multiple: true
  });

  return (
    <div>
      {/* Dropzone */}
      <div
        {...getRootProps()}
        style={{
          border: "2px dashed #aaa",
          padding: 20,
          cursor: "pointer",
          borderRadius: 8,
          textAlign: "center",
          width: "100%"
        }}
      >
        <input {...getInputProps()} />
        <p>Для загрузки файлов нажмите или перетащите их сюда</p>
      </div>

      {/* File list */}
      <div>
        {files.map(f => (
          <div
            key={f.id}
            style={{
              marginTop: 15,
              marginBottom: 15,
              padding: 12,
              border: "1px solid #ddd",
              borderRadius: 6,
              display: "flex",
              alignItems: "center",
              justifyContent: "space-between"
            }}
          >
            <div style={{ display: "flex", alignItems: "center", gap: 10, flex: 1 }}>
              <FileText size={22} />

              <div style={{ display: "flex", flexDirection: "column" }}>
                <div style={{ display: "flex", alignItems: "center", gap: 8 }}>
                  <strong>{f.file.name}</strong>
                </div>
              </div>
            </div>

            <div style={{ display: "flex", alignItems: "center", gap: 5 }}>
              <button
                onClick={() => removeFile(f.id)}
                style={{
                  border: "none",
                  background: "transparent",
                  cursor: "pointer",
                  padding: 5,
                  color: "#555",
                  outline: "none",
                  boxShadow: "none"
                }}
                onMouseOver={(e) => (e.currentTarget.style.color = "red")}
                onMouseOut={(e) => (e.currentTarget.style.color = "#555")}
              >
                <Trash2 size={18} />
              </button>
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}