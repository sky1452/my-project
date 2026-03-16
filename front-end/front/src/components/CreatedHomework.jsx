
export function CreatedHomework({disciplineId, group}){

return(

<div style={{marginTop:"30px"}}>

<h3>Созданные задания</h3>

<p>Дисциплина: {disciplineId}</p>
<p>Группа: {group}</p>

</div>

);

}