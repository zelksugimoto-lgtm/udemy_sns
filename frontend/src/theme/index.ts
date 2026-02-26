import { lightTheme } from './lightTheme';
import { darkTheme } from './darkTheme';
import { THEMES } from '../utils/constants';
import type { ThemeType } from '../utils/constants';

export { lightTheme, darkTheme };

export const getTheme = (themeType: ThemeType) => {
  switch (themeType) {
    case THEMES.LIGHT:
      return lightTheme;
    case THEMES.DARK:
      return darkTheme;
    case THEMES.CUSTOM1:
      // TODO: Implement custom theme 1
      return lightTheme;
    case THEMES.CUSTOM2:
      // TODO: Implement custom theme 2
      return lightTheme;
    default:
      return lightTheme;
  }
};
