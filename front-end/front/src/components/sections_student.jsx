import { Link } from 'react-router-dom';

const sections = [
  {
    name: 'Домашняя страница',
    id: 1,
    link: '/profile_student',
  },
  {
    name: 'Успеваемость',
    id: 2,
    link: '/progress_student',
  },
  {
    name: 'Расписание',
    id: 3,
    link: '/schedule_student',
  },
  {
    name: 'Предстоящие события',
    id: 4,
    link: '/events_student',
  },
  {
    name: 'Мои задания',
    id: 5,
    link: '/homework_student',
  },
  {
    name: 'Рейтинг',
    id: 6,
    link: '/rating_student',
  },
];
export function Sections_student() {
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