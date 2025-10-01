import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';

import { MainLayout_teacher } from './components/MainLayout_teacher';
import { Datap } from './components/profil_teacher';
import { SchedulePage } from './components/schedule_teacher';
import  Login  from './components/vhod';
import{ ProgressPage } from './components/progress_teacher'

import { MainLayout_student } from './components/MainLayout_student';
import { Datap_student } from './components/profil_student';
import { SchedulePage_student } from './components/schedule_student';
import{ ProgressPage_student } from './components/progress_student'

import './styles/layout.css';
import './styles/profil_teacher.css';
import './styles/stylet.css';
import './styles/Login.css';
import './styles/schedule_teacher.css';
import './styles/progress_teacher.css';
import './styles/progress_student.css';
function GroupsPage() {
  return <h2>Учебные группы и дисциплины</h2>;
}
function EventsPage() {
  return <h2>Предстоящие события</h2>;
}

function HomeworkPage() {
  return <h2>Отчёт о проверке домашнего задания</h2>;
}

function Rating_studentPage() {
  return <h2>Рейтинг студента</h2>;
}
function Events_studentPage() {
  return <h2>Предстоящие события</h2>;
}
function Homework_studentPage() {
  return <h2>Предоставление домашнего задания</h2>;
}
function App() {
  return (
    <Router>
      <Routes>
        {/* При заходе на "/" — редирект на "/login" */}
        <Route path="/" element={<Navigate to="/login" replace />} />

        <Route path="/login" element={<Login />} />

        {/* Все остальные маршруты с общей оболочкой */}
        <Route element={<MainLayout_teacher />}>
          <Route path="/profile_teacher" element={<Datap />} />
          <Route path="/groups_teacher" element={<GroupsPage />} />{/*костыль*/}
          <Route path="/schedule_teacher" element={<SchedulePage />} />
          <Route path="/events_teacher" element={<EventsPage />} />{/*костыль*/}
          <Route path="/progress_teacher" element={<ProgressPage />} />
          <Route path="/homework_teacher" element={<HomeworkPage />} />{/*костыль*/}
          </Route>
          <Route element={<MainLayout_student />}>
          <Route path="/profile_student" element={<Datap_student />} />
          <Route path="/rating_student" element={<Rating_studentPage />} />{/*костыль*/}
          <Route path="/schedule_student" element={<SchedulePage_student />} />
          <Route path="/events_student" element={<Events_studentPage />} />{/*костыль*/}
          <Route path="/progress_student" element={<ProgressPage_student />} />
          <Route path="/homework_student" element={<Homework_studentPage />} /> {/*костыль*/}
        </Route>
      </Routes>
    </Router>
  );
}

export default App;
