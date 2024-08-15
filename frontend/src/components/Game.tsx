"use client";

import { useState, useEffect, useRef } from "react";
import { Stage, Layer, Rect, Group, Circle } from "react-konva";

// Определяем интерфейсы для позиции игрока и смещения камеры

interface Position {
  x: number;
  y: number;
}

interface PlayerProps {
  position: Position;
}

const Player = ({ position }: PlayerProps) => {
  return (
    <Rect
      x={position.x * 50}
      y={position.y * 50}
      width={50}
      height={50}
      fill="blue"
    />
  );
};

const Grid = ({
  cellSize,
  cameraX,
  cameraY,
}: {
  cellSize: number;
  cameraX: number;
  cameraY: number;
}) => {
  const gridLines = [];

  // Устанавливаем диапазон линий, которые необходимо отрисовать на основе положения камеры
  const left = Math.floor(cameraX / cellSize) * cellSize;
  const right =
    left + cellSize * Math.ceil(window.innerWidth / cellSize) + cellSize;
  const top = Math.floor(cameraY / cellSize) * cellSize;
  const bottom =
    top + cellSize * Math.ceil(window.innerHeight / cellSize) + cellSize;

  // Рисуем вертикальные линии
  for (let x = left; x < right; x += cellSize) {
    gridLines.push(
      <Rect
        key={`vertical-${x}`}
        x={x}
        y={top}
        width={1}
        height={bottom - top}
        fill="lightgrey"
      />
    );
  }

  // Рисуем горизонтальные линии
  for (let y = top; y < bottom; y += cellSize) {
    gridLines.push(
      <Rect
        key={`horizontal-${y}`}
        x={left}
        y={y}
        width={right - left}
        height={1}
        fill="lightgrey"
      />
    );
  }

  return <>{gridLines}</>;
};

export default function Game() {
  const stageRef = useRef<any>(null);
  const [position, setPosition] = useState<Position>({ x: 0, y: 0 });
  const [pressedKeys, setPressedKeys] = useState<Set<string>>(new Set());
  const speed = 10; // Скорость движения игрока
  console.log(position);

  const updatePosition = (dx: number, dy: number) => {
    setPosition((prev) => {
      const newX = Math.min(window.innerWidth - 50, Math.max(0, prev.x + dx));
      const newY = Math.min(window.innerHeight - 50, Math.max(0, prev.y + dy));
      return { x: newX, y: newY };
    });
  };

  const movePlayer = () => {
    let dx = 0,
      dy = 0;
    if (pressedKeys.has("KeyW")) dy = -speed;
    if (pressedKeys.has("KeyS")) dy = speed;
    if (pressedKeys.has("KeyA")) dx = -speed;
    if (pressedKeys.has("KeyD")) dx = speed;

    // Обновляем позицию игрока
    if (dx !== 0 || dy !== 0) {
      const length = Math.sqrt(dx * dx + dy * dy);
      if (length > 0) {
        const normalizedX = ((dx / length) * speed) / 50;
        const normalizedY = ((dy / length) * speed) / 50;
        updatePosition(normalizedX, normalizedY);
      }
    }
  };

  useEffect(() => {
    const handleKeyDown = (event: KeyboardEvent) => {
      if (!pressedKeys.has(event.code)) {
        setPressedKeys((prev) => new Set(prev).add(event.code));
      }
    };

    const handleKeyUp = (event: KeyboardEvent) => {
      setPressedKeys((prev) => {
        const newSet = new Set(prev);
        newSet.delete(event.code);
        return newSet;
      });
    };

    window.addEventListener("keydown", handleKeyDown);
    window.addEventListener("keyup", handleKeyUp);

    const interval = setInterval(movePlayer, 1000 / 60);

    return () => {
      window.removeEventListener("keydown", handleKeyDown);
      window.removeEventListener("keyup", handleKeyUp);
      clearInterval(interval);
    };
  }, [pressedKeys]);

  useEffect(() => {
    if (stageRef.current) {
      const stage = stageRef.current;
      const width = stage.width();
      const height = stage.height();

      // Вычисляем новое положение группы с учетом позиции игрока
      const newGroupX = width / 2 - position.x * 50;
      const newGroupY = height / 2 - position.y * 50;

      stage.getLayers()[0].setAttr("x", newGroupX);
      stage.getLayers()[0].setAttr("y", newGroupY);
      stage.draw();
    }
  }, [position]);

  return (
    <Stage width={window.innerWidth} height={window.innerHeight} ref={stageRef}>
      <Layer>
        <Group>
          <Grid
            cellSize={50}
            cameraX={
              window.innerWidth / 2 + position.x * 50 - window.innerWidth
            }
            cameraY={
              window.innerHeight / 2 + position.y * 50 - window.innerWidth / 2
            }
          />
          <Circle x={500} y={525} radius={200} fill={"red"} />
          <Player position={position} />
        </Group>
      </Layer>
    </Stage>
  );
}
