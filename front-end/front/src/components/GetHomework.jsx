import { useQuery } from "@tanstack/react-query";
import { useParams, useNavigate } from "react-router-dom";
import { fetchHomeworks } from "./api";

export function ReceivedWorks({ disciplineId, group, userId }) {
  const navigate = useNavigate();
  const { disciplineSlug, homeworkId } = useParams();

  const { data: homeworks, isLoading, isError } = useQuery({
    refetchInterval: 1000,
    queryKey: ["homeworks", userId, group, disciplineId],
    queryFn: () => fetchHomeworks(userId, group, disciplineId),
    enabled: !!userId && !!group && !!disciplineId,
  });

  if (isLoading) return <p>Загрузка полученных работ...</p>;
  if (isError) return <p>Ошибка при загрузке полученных работ.</p>;

  return (
    <div style={{ marginTop: "30px" }}>
      Выберите созданное задание, чтобы посмотреть отправленные студентами работы:
      <div className="zadaniyes1">
        <div className="zadaniya-list1">
          {homeworks && homeworks.length > 0 ? (
            homeworks.map((hw) => (
              <div
                className="zadaniye1"
                key={hw.id}
                onClick={() =>
                  navigate(
                    `/homework_teacher/${disciplineId}/${disciplineSlug}/${encodeURIComponent(group)}/${hw.id}`
                  )
                }
                style={{
                  cursor: "pointer",
                  border:
                    String(hw.id) === String(homeworkId)
                      ? "2px solid #4f46e5"
                      : "",
                }}
              >
                {hw.title}
              </div>
            ))
          ) : (
            <div>У вас пока нет созданных заданий в этой дисциплине и группе.</div>
          )}
        </div>
      </div>
    </div>
  );
}