const PUBLIC_PAGE_THEME_CLASSES = ['dark', 'chartreuseDark'] as const;

export const isolatePublicPageTheme = () => {
  const root = document.documentElement;
  const activeClasses = PUBLIC_PAGE_THEME_CLASSES.filter((className) => root.classList.contains(className));
  root.classList.remove(...PUBLIC_PAGE_THEME_CLASSES);

  return () => {
    activeClasses.forEach((className) => root.classList.add(className));
  };
};
