import { CloseButton } from "@mantine/core";
import FileUploader from "./FileUploader";

export function HomeworkAnswerForm({
  bottomRef,
  isSubmitting,
  onClose,
  onFilesChange,
  onCommentChange,
  onSubmit,
}) {
  return (
    <div
      ref={bottomRef}
      className="my-homework2"
      style={{ position: "relative" }}
    >
      <CloseButton
        onClick={onClose}
        style={{ position: "absolute", top: 8, right: 8 }}
        disabled={isSubmitting}
      />

      <div className="div001">Ответ на задание</div>

      <div className="div002">
        <FileUploader onFilesChange={onFilesChange} />
      </div>

      <div className="div004">Комментарий:</div>
      <div className="div005">
        <div
          contentEditable={!isSubmitting}
          className="editable"
          onInput={(e) => onCommentChange(e.currentTarget.textContent)}
        />
      </div>

      <div
        className="div006"
        onClick={!isSubmitting ? onSubmit : undefined}
        style={{
          opacity: isSubmitting ? 0.6 : 1,
          pointerEvents: isSubmitting ? "none" : "auto",
        }}
      >
        {isSubmitting ? "Отправка..." : "Отправить"}
      </div>
    </div>
  );
}