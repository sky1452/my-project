import { useEffect, useRef, useState } from "react";
import { useMutation } from "@tanstack/react-query";
import { X } from "lucide-react";
import { updateHomework } from "./api";

export function EditHomework({
  homework,
  teacherId,
  disciplineId,
  group,
  onClose,
  onSuccess,
  tableRef,
}) {
  const [nameHomework, setNameHomework] = useState("");
  const [description, setDescription] = useState("");
  const [deadline, setDeadline] = useState("");
  const [maxScore, setMaxScore] = useState("");

  const formRef = useRef(null);
  const titleRef = useRef(null);
  const descriptionRef = useRef(null);
  const maxScoreRef = useRef(null);

  useEffect(() => {
    if (!homework) return;

    setNameHomework(homework.title || "");
    setDescription(homework.description || "");
    setDeadline(formatDateForInput(homework.deadline) || "");
    setMaxScore(String(homework.max_score ?? ""));

    if (titleRef.current) {
      titleRef.current.textContent = homework.title || "";
    }

    if (descriptionRef.current) {
      descriptionRef.current.textContent = homework.description || "";
    }

    if (maxScoreRef.current) {
      maxScoreRef.current.textContent = String(homework.max_score ?? "");
    }

    if (formRef.current) {
      formRef.current.scrollIntoView({ behavior: "smooth", block: "start" });
    }
  }, [homework]);

  const updateHomeworkMutation = useMutation({
    mutationFn: updateHomework,
    onSuccess: () => {
      onSuccess();
    },
  });

  const handleUpdateHomework = () => {
  const payload = {
    homeworkId: homework.id,
    teacherId,
    disciplineId,
    group,
    title: nameHomework,
    description,
    deadline,
    maxScore: Number(maxScore),
  };

  console.log("UPDATE PAYLOAD", payload);

  updateHomeworkMutation.mutate(payload);
};

  const handleClose = () => {
    onClose();

    setTimeout(() => {
      if (tableRef?.current) {
        tableRef.current.scrollIntoView({ behavior: "smooth", block: "start" });
      }
    }, 0);
  };

  return (
    <div ref={formRef} className="create-homework">
      <div className="edit-homework-close">
        <X size={20} onClick={handleClose} />
      </div>

      <div className="div1">Название задания:</div>

      <div className="div2">
        <div
          ref={titleRef}
          contentEditable={true}
          className="editable"
          onInput={(e) => setNameHomework(e.currentTarget.textContent || "")}
        ></div>
      </div>

      <div className="div3">Описание:</div>

      <div className="div4">
        <div
          ref={descriptionRef}
          contentEditable={true}
          className="editable"
          onInput={(e) => setDescription(e.currentTarget.textContent || "")}
        ></div>
      </div>

      <div className="div5">Дедлайн:</div>

      <input
        className="div6"
        type="datetime-local"
        value={deadline}
        onChange={(e) => setDeadline(e.target.value)}
      />

      <div className="div7">Максимальный балл:</div>

      <div className="div8">
        <div
          ref={maxScoreRef}
          contentEditable={true}
          className="editable"
          onBeforeInput={(e) => {
            const text = e.currentTarget.textContent || "";
            if (e.data && !/^\d$/.test(e.data)) {
              e.preventDefault();
              return;
            }
            if (text.length >= 3) {
              e.preventDefault();
            }
          }}
          onInput={(e) => {
            setMaxScore(e.currentTarget.textContent || "");
          }}
        ></div>
      </div>

      <div className="div11" onClick={handleUpdateHomework}>
        Сохранить изменения
      </div>

      {updateHomeworkMutation.isPending && <p>Сохранение изменений...</p>}
      {updateHomeworkMutation.isError && (
        <p>Ошибка: {updateHomeworkMutation.error.message}</p>
      )}
    </div>
  );
}

function formatDateForInput(dateString) {
  if (!dateString) return "";

  const date = new Date(dateString);

  if (Number.isNaN(date.getTime())) return "";

  const year = date.getFullYear();
  const month = String(date.getMonth() + 1).padStart(2, "0");
  const day = String(date.getDate()).padStart(2, "0");
  const hours = String(date.getHours()).padStart(2, "0");
  const minutes = String(date.getMinutes()).padStart(2, "0");

  return `${year}-${month}-${day}T${hours}:${minutes}`;
}