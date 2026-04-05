export const fetchDiscipline = async (disciplineId) => {
  const res = await fetch(`http://localhost:8081/discipline/${disciplineId}`);
  return res.json();
};

export const fetchStudentTasks = async (disciplineId, userId) => {
  const res = await fetch(
    `http://localhost:8081/discipline/${disciplineId}/student/${userId}`
  );
  return res.json();
};
export const fetchProgress = async (userId, disciplineId) => {
  const res = await fetch(
    `http://localhost:8081/progress/${userId}/${disciplineId}`
  );
  return res.json();
};

export const fetchTaskById = async (taskId) => {
  const res = await fetch(`http://localhost:8081/tasks/${taskId}`);
  return res.json();
};
export async function submitHomework(taskId, comment, files, userId, disciplineId) {
  const formData = new FormData();

  formData.append("task_id", taskId);
  formData.append("student_id", userId); //временно
  formData.append("discipline_id", disciplineId);
  formData.append("comment", comment);

  files.forEach((file) => {
    formData.append("files", file);
  });

  return fetch("http://localhost:8081/submissions", {
    method: "POST",
    body: formData,
  });
}
export const fetchFileById = async (taskId, userId) => {
  const res = await fetch(`http://localhost:8081/tasks/${taskId}/student/${userId}/files`);
  return res.json();
};

export const fetchHomeworks = async (userId, group, disciplineId) => {
  const res = await fetch(
    `http://localhost:8081/getHomeworks/${userId}/${group}/${disciplineId}`
  );

  return res.json();
};
export async function fetchHomeworkById(userId, group, disciplineId, homeworkId) {
  const response = await fetch(
    `http://localhost:8081/getHomework/${userId}/${encodeURIComponent(group)}/${disciplineId}/${homeworkId}`
  );

  if (!response.ok) {
    throw new Error("Ошибка при загрузке задания");
  }

  return response.json();
}
export async function fetchHomeworkSubmissions(taskId, group) {
  const res = await fetch(
    `http://localhost:8081/homeworks/${taskId}/group/${encodeURIComponent(group)}/submissions`
  );

  if (!res.ok) {
    throw new Error("Ошибка загрузки submissions");
  }

  return res.json();
}
export async function updateSubmissionScore(submissionId, score) {
  const res = await fetch(
    `http://localhost:8081/submissions/${submissionId}/score`,
    {
      method: "PUT",
      headers: {
        "Content-Type": "application/json",
      },
      body: JSON.stringify({ score }),
    }
  );

  if (!res.ok) {
    const text = await res.text();
    console.error("Ошибка при обновлении оценки:", text);
    throw new Error(text || "Ошибка обновления оценки");
  }

  return res.json();
}
export async function fetchSubmissionScore(taskId, userId) {
  const res = await fetch(
    `http://localhost:8081/tasks/${taskId}/student/${userId}/score`
  );

  if (!res.ok) {
    throw new Error("Ошибка загрузки оценки");
  }

  return res.json();
}
export async function updateHomeworkAnswer({
  taskId,
  userId,
  disciplineId,
  comment,
  newFiles,
  keptFileIndexes,
}) {
  const formData = new FormData();

  formData.append("comment", comment || "");
  formData.append("disciplineId", disciplineId);
  formData.append("keptFileIndexes", JSON.stringify(keptFileIndexes || []));

  (newFiles || []).forEach((file) => {
    formData.append("files", file);
  });

  const res = await fetch(
    `http://localhost:8081/tasks/${taskId}/student/${userId}/submission`,
    {
      method: "PUT",
      body: formData,
    }
  );

  if (!res.ok) {
    const text = await res.text();
    throw new Error(text || "Ошибка редактирования ответа");
  }

  return res.json();
}