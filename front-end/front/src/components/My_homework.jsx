import { useState, useRef, useEffect } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import {
  fetchTaskById,
  fetchDiscipline,
  submitHomework,
  fetchFileById,
} from "./api";
import { HomeworkInfo } from "./HomeworkInfo";
import { HomeworkAnswerForm } from "./HomeworkAnswerForm";

export function MyHomeworks() {
  const user = JSON.parse(localStorage.getItem("user"));
  const userId = user?.userId;

  const bottomRef = useRef(null);
  const [uploadedFiles, setUploadedFiles] = useState([]);
  const [comment, setComment] = useState("");
  const [createAnswer, setCreateAnswer] = useState(false);
  const [isSubmitting, setIsSubmitting] = useState(false);

  const { disciplineId, taskId } = useParams();
  const queryClient = useQueryClient();

  const { data: submissionData, isLoading: submissionLoading } = useQuery({
    queryKey: ["submission-files", taskId, userId],
    queryFn: () => fetchFileById(taskId, userId),
    enabled: !!taskId && !!userId,
    refetchInterval: 1000,
  });

  const { data: task, isLoading: taskLoading } = useQuery({
    queryKey: ["task", taskId],
    queryFn: () => fetchTaskById(taskId),
    refetchInterval: 1000,
    initialData: () => {
      const allTasks = queryClient.getQueryData(["tasks"]);
      return allTasks?.tasks?.find((t) => t.id === Number(taskId));
    },
  });

  const { data: discipline, isLoading: disciplineLoading } = useQuery({
    queryKey: ["discipline", disciplineId],
    queryFn: () => fetchDiscipline(disciplineId),
    refetchInterval: 1000,
    initialData: () => {
      const allDisciplines = queryClient.getQueryData(["disciplines"]);
      return allDisciplines?.disciplines?.find(
        (d) => d.id === Number(disciplineId)
      );
    },
  });

  useEffect(() => {
    if (createAnswer) {
      bottomRef.current?.scrollIntoView({
        behavior: "smooth",
        block: "start",
      });
    }
  }, [createAnswer]);

  const handleSubmit = async () => {
    try {
      setIsSubmitting(true);

      await submitHomework(
        taskId,
        comment,
        uploadedFiles,
        userId,
        disciplineId
      );

      alert("Отправлено!");
      setCreateAnswer(false);
      setComment("");
      setUploadedFiles([]);
    } catch (e) {
      alert("Ошибка отправки");
    } finally {
      setIsSubmitting(false);
    }
  };

  if (disciplineLoading) return <div>Загрузка дисциплины...</div>;
  if (taskLoading) return <div>Загрузка задания...</div>;
  if (!task) return <div>Задание не найдено</div>;

  return (
    <div className={`progresst ${isSubmitting ? "blurred-content" : ""}`}>
      <div>{discipline?.name}</div>

      <HomeworkInfo
        task={task}
        taskId={taskId}
        userId={userId}
        submissionData={submissionData}
        submissionLoading={submissionLoading}
        createAnswer={createAnswer}
        onCreateAnswer={() => setCreateAnswer(true)}
      />

      {createAnswer && (
        <HomeworkAnswerForm
          bottomRef={bottomRef}
          isSubmitting={isSubmitting}
          onClose={() => setCreateAnswer(false)}
          onFilesChange={setUploadedFiles}
          onCommentChange={setComment}
          onSubmit={handleSubmit}
        />
      )}
    </div>
  );
}