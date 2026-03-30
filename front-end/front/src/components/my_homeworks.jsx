import  { useState } from "react";
import { useQuery, useQueryClient } from "@tanstack/react-query";
import { useParams } from "react-router-dom";
import { fetchTaskById } from "./api";
import { fetchDiscipline } from "./api";
import { CloseButton } from "@mantine/core";
import { useUploady } from "@rpldy/uploady";


export function MyHomeworks() {
    const { upload } = useUploady();
    const {disciplineId} = useParams();
    const { taskId } = useParams(); 
    const queryClient = useQueryClient();
    const [createAnswer, setCreateAnswer] = useState(false);
    const { data: task, isLoading } = useQuery({
    queryKey: ["task", taskId],
    queryFn: () => fetchTaskById(taskId),
    refetchInterval: 1000,
    initialData: () => {
      const allTasks = queryClient.getQueryData(["tasks"]);
      return allTasks?.tasks?.find(t => t.id === Number(taskId));
    },
  });
    const { data: discipline, Loading } = useQuery({
    queryKey: ["discipline", disciplineId],
    queryFn: () => fetchDiscipline(disciplineId),
    refetchInterval: 1000,

    initialData: () => {
      const allDisciplines = queryClient.getQueryData(["disciplines"]);
      return allDisciplines?.disciplines?.find(d => d.id === Number(disciplineId));
    },
  });

  if (Loading) return <div>Загрузка дисциплины...</div>;
  if (isLoading) return <div>Загрузка задания...</div>;
  if (!task) return <div>Задание не найдено</div>;

  return (
    <div className="progresst">
      <div>{discipline?.name}</div>
      <div className="my-homework1">
      <div className="div01"> {task.title} </div>
        <div className="div02">Описание:  </div>
        <div className="div03"> {task.description} 123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890123456789012345678901234567890 </div>
          <div className="div04">Максимальная оценка: </div>
              <div className="div05"> {task.max_score} </div>
             <div className="div06">Задание создано: </div>
            <div className="div07"> {task.created_at} </div>
              <div className="div08">Крайний срок сдачи: </div>
                <div className="div09"> {task.deadline} </div>
              <div className="div010">Ответ в виде файла: </div>
              <div className="div011">- </div>
               <div className="div012">Комментарий: </div>
                <div className="div013">- </div>
                <div className="div014">Последнее изменение: </div>
                <div className="div015">(надо подставить дату)-</div>
                <div className="div016"> Состояние оценивания:</div>
                <div className="div017">- </div>
                <div className="div018" onClick={() => setCreateAnswer(true)}>Создать ответ</div>
        </div>
       {createAnswer && (
  <div className="my-homework2" style={{ position: "relative" }}>
    
    <CloseButton
      onClick={() => setCreateAnswer(false)}
      style={{ position: "absolute", top: 8, right: 8 }}
    />

    <div className="div001">Ответ на задание</div>
    <div className="div002">Прикрепить файл:</div>
    <div className="div003">3</div>
    <div className="div004">Комментарий:</div>
    <div className="div005">5</div>
    <div className="div006">Отправить</div>

  </div>
)}
      </div>
    
  );
}
