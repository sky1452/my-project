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
export const fetchProgress = async (userId) => {
  const res = await fetch(
    `http://localhost:8081/progress/${userId}`
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