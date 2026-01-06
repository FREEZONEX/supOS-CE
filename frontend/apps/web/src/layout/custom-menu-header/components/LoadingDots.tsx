const LoadingDots = ({ color = '#333', size = '5px', speed = 0.8 }) => {
  const dotStyle = {
    backgroundColor: color,
    width: size,
    height: size,
    borderRadius: '50%',
    animation: `bounce ${speed}s ease-in-out infinite`,
  };

  return (
    <div style={{ display: 'flex', alignItems: 'center', gap: size, height: '100%' }}>
      <div style={{ ...dotStyle, animationDelay: '0s' }} />
      <div style={{ ...dotStyle, animationDelay: `${speed * 0.33}s` }} />
      <div style={{ ...dotStyle, animationDelay: `${speed * 0.66}s` }} />

      <style>{`
        @keyframes bounce {
          0%, 100% { transform: translateY(0); }
          50% { transform: translateY(-${parseInt(size) * 1.5}px); }
        }
      `}</style>
    </div>
  );
};

export default LoadingDots;
