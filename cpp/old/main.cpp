#include <SDL2/SDL.h>
#include <SDL2/SDL_image.h>
#include <SDL2/SDL_ttf.h>
#include <iostream>
#include <vector>

#define WIDTH 1600
#define HEIGHT 900
#define FONT_SIZE 32
#define TILE_SIZE 64
#define MAP_SIZE 17
#define SPEED 17

#define PI 3.14159265359

struct object {
 SDL_Rect dest, src;
 int img;
 SDL_Color c;
 int x, y;
 double angle;
 float vel;
 int tick;
} car, center;

int grid[MAP_SIZE][MAP_SIZE];
std::vector<SDL_Texture*> image;
int i1, i2;

bool running;

SDL_Renderer* renderer;
SDL_Window* window;
TTF_Font *font;
SDL_Color fcolor;

SDL_Point mouse;
int frameCount, timerFPS, lastFrame, fps;

bool up, down, right, left;
int x, y;

int setImage(std::string filename) {
 image.push_back(IMG_LoadTexture(renderer, filename.c_str()));
 return image.size()-1;
}

void draw(object o) {
 SDL_RenderCopyEx(renderer, image[o.img], &o.src, &o.dest, o.angle, NULL, SDL_FLIP_NONE);
}

void drawRect(SDL_Rect r) {
 SDL_RenderFillRect(renderer, &r);
}

void drawOutline(SDL_Rect r) {
 SDL_RenderDrawRect(renderer, &r);
}

SDL_Surface *surface;
SDL_Texture *texture;
SDL_Rect wrect;
void write(std::string text, int x, int y) {
//TTF_SetFontOutline(font, 1);
 if (font == NULL) {
  fprintf(stderr, "error: font not found\n");
  exit(EXIT_FAILURE);
 }
 fcolor.r = 0;
 fcolor.g = 0;
 fcolor.b = 0;
 const char* t = text.c_str();
 surface = TTF_RenderText_Solid(font, t, fcolor);
 texture = SDL_CreateTextureFromSurface(renderer, surface);
 wrect.w = surface->w;
 wrect.h = surface->h;
 wrect.x = x-wrect.w;
 wrect.y = y-wrect.h;
 SDL_FreeSurface(surface);
 SDL_RenderCopy(renderer, texture, NULL, &wrect);
 SDL_DestroyTexture(texture);
}

const float halfC = PI / 180;
float get_degrees(float input) {
    return input * halfC;
}

void camera(int ox, int oy) {
 x=ox + WIDTH/2;
 y=oy + HEIGHT/2;
}
void camera(object o) {camera(o.x,o.y);}

int dx, dy;
void update() {
 camera(car);
 //camera(center);
 car.dest.x=car.x;
 car.dest.y=car.y;
 //camera((MAP_SIZE/2)*TILE_SIZE, (MAP_SIZE/2)*TILE_SIZE);
 dx = cos(get_degrees(car.angle))*car.vel;
 dy = sin(get_degrees(car.angle))*car.vel;
 //if(down || up) {car.y+=dy;car.x+=dx;}//SPEED;
 car.y+=dy;car.x+=dx;//SPEED;
 if(up) car.vel+=0.5;
 if(down) car.vel-=0.5;
 if(car.vel>2 || car.vel<-2) {
  if(left) car.angle-=2;//car.x+=SPEED;
  if(right) car.angle+=2;//x-=SPEED;
 }
 if(up || down) car.tick++;
 if(!up && !down) {
  if(car.vel>0) car.vel-=0.5;
  if(car.vel<0) car.vel+=0.5;
 }
 if(car.vel==0)car.tick=0;
 if(car.vel>SPEED/2 && car.tick<100) car.vel=SPEED/2;
 if(car.vel<-SPEED/2 && car.tick<100) car.vel=-SPEED/2;
 if(car.vel>SPEED) car.vel=SPEED;
 if(car.vel<-SPEED) car.vel=-SPEED;
}

const Uint8 *keystates;
void input() {
    left=right=down=up=0;
    SDL_Event e;
    keystates = SDL_GetKeyboardState(NULL);
    while(SDL_PollEvent(&e)) {
        if(e.type == SDL_QUIT) running=false;
    }
    if(keystates[SDL_SCANCODE_ESCAPE]) running=false;
    if(keystates[SDL_SCANCODE_W]) up=1;
    if(keystates[SDL_SCANCODE_S]) down=1;
    if(keystates[SDL_SCANCODE_A]) left=1;
    if(keystates[SDL_SCANCODE_D]) right=1;

    SDL_GetMouseState(&mouse.x, &mouse.y);
}

object d;
void render() {
 SDL_SetRenderDrawColor(renderer, 102, 75, 71, 255);
 SDL_RenderClear(renderer);
 frameCount++;
 int timerFPS = SDL_GetTicks()-lastFrame;

 d.dest.w=d.dest.h=TILE_SIZE;
 d.src.w=d.src.h=TILE_SIZE;
 d.src.x=d.src.y=0;
 /*int sx=x/TILE_SIZE - 1;
 int sy=y/TILE_SIZE - 1;
 if(sx<0) sx=0;
 if(sy<0) sy=0;*/
 for(int i=0; i<MAP_SIZE; i++) {
  for(int j=0; j<MAP_SIZE; j++) {
   d.dest.x=i*TILE_SIZE-x;
   d.dest.y=j*TILE_SIZE-y;
   if(d.dest.x>WIDTH) {break;}
   if(d.dest.x+TILE_SIZE<0) {break;}
   //if(d.dest.y>HEIGHT) {break;}
   //if(d.dest.y+TILE_SIZE<0) {break;}
   if(d.dest.x+TILE_SIZE>0 && d.dest.x-TILE_SIZE<WIDTH && d.dest.y+TILE_SIZE>0 && d.dest.y-TILE_SIZE<HEIGHT) {
   //if(grid[i][j]==2) //SDL_SetRenderDrawColor(renderer, 240, 0, 240, 255);
   //if(grid[i][j]==1) //SDL_SetRenderDrawColor(renderer, 240, 240, 0, 255);
    d.img=grid[i][j];
    draw(d);
    SDL_SetRenderDrawColor(renderer, 244, 147, 94, 255);
    drawOutline(d.dest);
   }
  }
 }

 d.dest.x=(MAP_SIZE/2)*TILE_SIZE-x;
 d.dest.y=(MAP_SIZE/2)*TILE_SIZE-y;
// d.x=d.dest.x; d.y=d.dest.y;
 SDL_SetRenderDrawColor(renderer, 0, 255, 0, 255);
 drawOutline(d.dest);
 center=d;
 draw(car);
 write(std::to_string(car.x) + ", " + std::to_string(car.y) + " " + std::to_string(car.vel), mouse.x, mouse.y);

 SDL_RenderPresent(renderer);
}

void init() {
 bool sw=1;
 i1=setImage("1.png");
 i2=setImage("2.png");
 for(int i=0; i<MAP_SIZE; i++) {
  for(int j=0; j<MAP_SIZE; j++) {
   grid[j][i]=i1;
   if(sw)grid[j][i]=i2;
   sw=!sw;
  }
 }
 x=550;y=250;
 car.src.x=car.src.y=0;
 car.src.w=12;car.src.h=7;
 car.dest.w=33; car.dest.h=20;
 car.dest.x=WIDTH/2;
 car.dest.y=HEIGHT/2;
 car.img=setImage("car.png");
 car.x=car.y=300;
 car.vel=0;
}

int main() {
    running=1;
    static int lastTime=0;
    SDL_SetHint(SDL_HINT_RENDER_SCALE_QUALITY, "0");
    if(SDL_Init(SDL_INIT_EVERYTHING) < 0) std::cout << "Failed at SDL_Init()" << std::endl;
    //if(SDL_CreateWindowAndRenderer(WIDTH, HEIGHT, SDL_WINDOW_FULLSCREEN, &window, &renderer) < 0) std::cout << "Failed at SDL_CreateWindowAndRenderer()" << std::endl;
    window = SDL_CreateWindow("Game", SDL_WINDOWPOS_UNDEFINED, SDL_WINDOWPOS_UNDEFINED, WIDTH, HEIGHT, SDL_WINDOW_SHOWN);
    renderer = SDL_CreateRenderer(window, -1, SDL_RENDERER_ACCELERATED | SDL_RENDERER_PRESENTVSYNC);
    SDL_SetWindowFullscreen(window, SDL_WINDOW_FULLSCREEN);
    TTF_Init();
    font = TTF_OpenFont("pricedown.ttf", FONT_SIZE);
    if(font == NULL) std::cout << "failed to load font" << std::endl;

    init();

    while(running) {
        lastFrame=SDL_GetTicks();
        if(lastFrame>=(lastTime+1000)) {
            lastTime=lastFrame;
            fps=frameCount;
            frameCount=0;
        }

        update();
        input();
        render();
    }
    TTF_CloseFont(font);
    SDL_DestroyRenderer(renderer);
    SDL_DestroyWindow(window);
    SDL_Quit();
}
