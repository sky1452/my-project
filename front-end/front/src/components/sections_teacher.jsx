import { Link } from 'react-router-dom';

const sections = [
  {
    name: 'Домашняя страница',
    id: 1,
    link: '/profile_teacher',
  },
  {
    name: 'Группы и дисциплины',
    id: 2,
    link: '/groups_teacher',
  },
  {
    name: 'Текущее расписание',
    id: 3,
    link: '/schedule_teacher',
  },
  {
    name: 'Предстоящие события',
    id: 4,
    link: '/events_teacher',
  },
  {
    name: 'Формирование успеваемости',
    id: 5,
    link: '/progress_teacher',
  },
  {
    name: 'Проверка заданий',
    id: 6,
    link: '/homework_teacher',
  },
];

export function SectionsTeacher() {
  return (
    <div className="sections-container">
      {sections.map((section) => (
        <Link
          key={section.id}
          className="sections"
          to={section.link}
        >
          <p>{section.name}</p>
        </Link>
      ))}
    </div>
  );
}
